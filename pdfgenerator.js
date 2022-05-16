// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

const fs = require('fs');
const path = require('path');
const util = require('util')
const Handlebars = require('handlebars');
const puppeteer = require('puppeteer');

const data = require('./output.json')

function commaFormat(number) {
	return (number/100).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function creditCommaFormat(number) {
	return (number/-100).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}


(async function() {
  try {

    const browser = await puppeteer.launch();
    const page = await browser.newPage();

    const templatePath = path.resolve('financials.html')
    var contents = fs.readFileSync(templatePath, 'utf8');

		Handlebars.registerHelper('commaFormat', function(number) {
			return commaFormat(number);
		});

		Handlebars.registerHelper('creditCommaFormat', function(number) {
			return creditCommaFormat(number);
		});

    Handlebars.registerHelper('if_eq', function(a, b, opts) {
				if (a == b) {
						return opts.fn(this);
				} else {
						return opts.inverse(this);
				}
    });

		Handlebars.registerHelper("debug", function(optionalValue) {
			console.log("Current Context");
			console.log("====================");
			console.log(this);

			if (optionalValue) {
				console.log("Value");
				console.log("====================");
				console.log(optionalValue);
			}
		});

    const template = Handlebars.compile(contents)
    await page.setContent(template({data}))

    await page.emulateMedia('screen');
    await page.pdf({
      path: 'mypdf.pdf',
      format: 'a4',
      printBackground: true
    });

    console.log('done');
    browser.close();
    process.exit();

  } catch (e) {
    console.log(e);
  }
})();