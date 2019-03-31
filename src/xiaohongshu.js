const spider = require('../lib/spider');
const reset = require('../lib/resetChorme');
const log = require("../lib/log");
const MongoClient = require('mongodb').MongoClient;
const mongo_dsn = "mongodb://localhost:27017/";
const trim = require('trim');
const mongoConnect  =  (mongo_dsn)=>{
    return new Promise((resolve, reject) =>{
        MongoClient.connect(mongo_dsn, { useNewUrlParser: true }, function(err, db) {
            if (err) reject(err);
            resolve(db);
        });
    })
}

/**id
 * 获取dom元素内的文本
 */
async function getElementTitle(page) {
    try {
        return await page.evaluate(() => {
            let text = '';
            if (document.querySelector('body > div.framework > div > div.card-note.pc-container > div.left-card > div.note-top > h1')) {
                text = document.querySelector('body > div.framework > div > div.card-note.pc-container > div.left-card > div.note-top > h1').textContent;
            }else if (document.querySelector('body > div.framework > div > div.card-note.pc-container > div.left-card > div:nth-child(3) > div > h1')){
                text = document.querySelector('body > div.framework > div > div.card-note.pc-container > div.left-card > div:nth-child(3) > div > h1').textContent;
            }
            return text;
        }, {});
    } catch (e) {
        return "";
    }
}

 function getElementUid(str) {
    const path = require('path');
    return path.basename(str);
}

/**id
 * 获取dom元素内的文本
 */
async function getElementText(page, selector) {
    try {
        await  page.waitForSelector(selector);
        return trim(await page.$eval(selector, ele => ele.textContent));
    } catch (e) {
        return "";
    }
}

/**
 * 获取dom元素内的html
 */
async function getElementHtml(page, selector) {
    try {
        await page.waitForSelector(selector);
        return (await page.$eval(selector, ele => ele.innerHTML));
    } catch (e) {
        return "";
    }
}

async function getElementHref(page, selector) {
    try {
        return await page.$eval(selector, ele => ele.href);
    } catch (e) {
        return "";
    }
}
/**
 * 获取分类列表
 * @param page
 * @param url
 * @returns {Promise<boolean>}
 */
async function getCategories(page,url) {
    try {
        await page.goto(url);
        await page.waitForSelector('body .note-tab');
        return await page.$$eval("body .note-tab a",urls=>{
            return urls.map((url)=>{
                return url.href;
            })
        })
    }catch (e) {
        return false;
    }
}

/**
 * 获取用户列表
 * @param page
 * @param url
 * @returns {Promise<boolean>}
 */
async function getUsersAndTasks(page,url) {
    try {
        await page.goto(url);
        await page.waitForSelector('body .note-box');
         let users =  await page.$$eval(".note-box .note-handle a",urls=>{
            return urls.map((url)=>{
                return url.href;
            });
        });
         let tasks = await page.$$eval(".note-box .note-info a",urls=>{
             return urls.map((url)=>{
                 return url.href;
             });
         });
         return {users,tasks}
    }catch (e) {
        return false;
    }
}


/**
 * 获取用户的发帖列表
 *
 * @param page
 * @param url
 * @returns {Promise<boolean>}
 */
async function getUserProfile(page,url) {
    try {
        await page.goto(url);
        await page.waitForSelector('body .notes-item');
        return await page.$$eval("body .each-list .note-info a",urls=>{
            return urls.map((url)=> url.href)
        })
    }catch (e) {
        return false;
    }
}


async  function getElementSrc(page,selector) {
    try {
        return trim(await page.$eval(selector, ele => ele.src));
    } catch (e) {
        return "";
    }
}

/**
 * 抓取内容界面内容
 * @param page
 * @param link
 * @param type
 * @returns {Promise<*>}
 */
async function runItem(page, link,type=1) {
    let data = {
        uid:"",
        type:type,
        item_uid:"",
        url:"",
        title:"",
        content:"",
        like:"",
        comment:"",
        star:"",
        commentinfo:"",
        pic:[],
    };
    let item_uid = getElementUid(link);
    try {
        await page.goto(link, {
            waitUntil: "domcontentloaded",
            timeout: 60000
        });
        data.item_uid = item_uid;
        data.type = type;
        data.url = link;

        const hasVideo = await page.evaluate( () => {
            return !!document.querySelector('video')
        });
        if(!hasVideo) {
            await page.waitForSelector('.carousel > ul');
            try {
                data.pic  = await page.$$eval(`.carousel > ul > li span`,pics=>{
                    return pics.map((pic)=>"https://" + pic.style.backgroundImage.match(/"\/\/(\S*)"/)[1])
                });
            }catch (e) {
                data.pic = [];
            }
        }else{
            data.type = 2;
            data.video_url = await getElementSrc(page,'video');
        }
        let user_url = await getElementHref(page,'.author-info');
        if(user_url) {
            data.uid = getElementUid(user_url);
        }
        data.title = await getElementTitle(page);
        data.content = await getElementHtml(page, '.content');
        data.like = await getElementText(page, '.like');
        data.comment = await getElementText(page, '.comment');
        data.star = await getElementText(page, '.star');
        return data;
    } catch (e) {
        log.error(e);
        return false;
    }
}
/**
 * Main函数
 */
(async ()=>{

    try {
        const browser = await spider.launch();
        const page = await browser.newPage();
        await page.setUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36");
        //页面响应事件
        await reset(page);
        const categories = await getCategories(page,"https://www.xiaohongshu.com/explore");
        const sleep = require('sleep-anywhere');
        const random = require('random');
        let v = random.int(5, 10);
        if(false === categories) {
           process.exit(0);
        }
        let tasks = [];
        let users = [];
        for (let i = 0; i < categories.length; i++) {
            let categoryUrl = categories[i];
            let data = await getUsersAndTasks(page,categoryUrl);
            if(data !== false) {
                users = [...users,...data.users];
                tasks = [...tasks,...data.tasks]
            }
        }
        users = new Set([...users]);

        for (let user of users) {
            let userTask = await getUserProfile(page,user);
            if(false !== userTask) {
                tasks = [...tasks,...userTask]
            }
        }
        const fs = require('fs');
        //爬取所有的列表
        tasks = new Set([...tasks]);
        const mongo = await mongoConnect(mongo_dsn);
        const dbo = mongo.db("xiaohongshu");
        for (let task of tasks) {
            const data = await runItem(page,task,1);
            if(false === data) {
                continue;
            }
            try {
                dbo.collection("items").insertOne(data, function(err, res) {
                    if (err) return err;
                    console.log(task,"文档插入成功");
                });
            }catch (e) {}

        }
        await browser.close();
        process.exit();
    } catch (e) {
        log.error(e);
        process.exit(-1);
    }
})();
