const spider = require('../lib/spider');
const axios = require("axios")
const Url = require('url-parse');
function extractItems() {
    const extractedElements = document.querySelectorAll('#main li.item a');
    const items = [];
    for (let element of extractedElements) {
      items.push(element.href);
    }
    return items;
  }
  
async function scrapeInfiniteScrollItems(
    page,
    extractItems,
    itemTargetCount,
    scrollDelay = 1000,
  ) {
    let items = [];
    try {
      let previousHeight;
      while (items.length < itemTargetCount) {
        items = await page.evaluate(extractItems);
        previousHeight = await page.evaluate('document.body.scrollHeight');
        await page.evaluate('window.scrollTo(0, document.body.scrollHeight)');
        await page.waitForFunction(`document.body.scrollHeight > ${previousHeight}`);
        await page.waitFor(scrollDelay);
      }
    } catch(e) { 
        console.log(e)
        return false;
    }
    return items;
}
async function getElementText(page, selector) {
    try {
        await  page.waitForSelector(selector);
        return await page.$eval(selector, ele => ele.textContent);
    } catch (e) {
        return e;
    }
}

function parseSalary(salary) {
    let temp = {
        min: 0,
        max: 0,
        avg: 0
    }
    //s = "13-22K·13薪"
    salary = salary || ""
    if (salary.length == 0) {
        return temp
    }
    t = salary.split("·")[0].split("-")
    if (t.length < 2) {
        return temp
    }
    let [min,max] = t
    step = 1
    if (max.indexOf("万") > 0 || max.indexOf("W") > 0) {
        step = 10000
    } else if (max.indexOf("千") > 0 || max.indexOf("K") > 0) {
        step = 1000
    }
    temp.min = parseFloat(min) * step
    temp.max = parseFloat(max) * step
    temp.avg = parseInt ((temp.min + temp.max) / 2)
    return temp
}
async function gotoDetail(page ,url) {
    try {
        await page.goto(url)
        console.log(url)
        await page.waitFor("body")
        let uParse = new Url(url)
        let path = uParse.pathname
        pathArr = path.split("/")
        postion_id = pathArr[pathArr.length - 1].replace(".html", "")
        let title = await getElementText(page ,"#main h1")
        let salary = await getElementText(page,"h1 .salary")
        position_name = title.replace(salary, "")
        let create = await  getElementText(page,".time")
        create_time = create.replace("更新于：", "")
        let vline = await page.$eval("#main > div.job-banner > p",e=>e.innerHTML)
        let vlineArr = vline.split(`<em class="vline"></em>`)
        let work_year = vlineArr[1]
        let educational = vlineArr[2]
        salary =  parseSalary(salary)
        let company_name = await getElementText(page, ".business-info h4")
        let body = await page.$$eval(".job-sec .text", e => e.map(l => l.textContent))
        let address = await getElementText(page, ".location-address")
        return {
            address,
            salary,
            create_time,
            body,
            company_name,
            postion_id,
            position_name,
            educational,
            work_year,
        }
        
    } catch (error) {
        return error
    }
}
(async ()=>{
    try {
        const browser = await spider.launch();
        const page = await browser.newPage();
        await page.setUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36");
        await page.setViewport({ width: 1280, height: 926 });
        await page.goto("https://m.zhipin.com/c101020100-p100103/?ka=position-100103")
        const items = await scrapeInfiniteScrollItems(page, extractItems, 100);
        let results = []
        for (let index = 0; index < items.length; index++) {
            let url = items[index];
            let res = await gotoDetail(page,url)
            if (false ===  res) {
                continue;
            } 
            results.push(res)
            await page.waitFor(2000);
        }
        axios.post('/api/', results)
            .then(function (response) {
                
            })
            .catch(function (error) {
                console.log(error);
            });
        await browser.close();
    }catch(e){
        console.log(e)
    }
})();