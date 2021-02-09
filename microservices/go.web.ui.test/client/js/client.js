
/* 
 * The MIT License (MIT)
 * 
 * Copyright 2020 Eugenio Parodi <eugenio.parodi.78@gmail.com>.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
"use strict";

var firstKeyTime = null;
var sessionId = null;

$(function(){
    sessionId = 
        (100000+   Math.floor(Math.random()*899999)) + "-" +
        (100000+   Math.floor(Math.random()*899999)) + "-" +
        (100000000+Math.floor(Math.random()*899999999)) ;

    /* I assume only one window resize as defined in the requirements */
    pushWindowSize();
    $(window).on('resize', pushWindowSize);

    /* Bind the input methods to capture the copy/paste and keypress events */
    $('input').each(function(){
      $(this).bind('keypress', function() { keyPressed(); });
      $(this).bind('copy', function() { pushCopyAndPaste(this.id, false); });
      $(this).bind('paste',function() { pushCopyAndPaste(this.id, true);  });	
      $(this).bind('cut',  function() { pushCopyAndPaste(this.id, false); });
    });

    $('form').submit(function( event ) {
        let keyPressTime = firstKeyTime ? $.now() - firstKeyTime : -1 ;
        console.log("Form submitted, time from the first keypress:",keyPressTime);
        post({
            "eventType": "timeTaken",
            "timeTaken": keyPressTime
        });
        event.preventDefault();
    });
});

function keyPressed(){
    console.log("[PUSH] keyPressed");
    if (!firstKeyTime){
        firstKeyTime = $.now();
    }
}

function pushCopyAndPaste(id, val){
    console.log("[PUSH] CopyandPaste id:",id," val:",val);
    post({
        "eventType": "copyAndPaste",
        "pasted": val,
        "formId": id
      });
}

function pushWindowSize(){
    /* disable the event to avoid unused calls */
    $(window).off('resize', pushWindowSize);
    let w = $(window).width();
    let h = $(window).height();
    console.log("[PUSH] size w:",w," h:",h);
    post({
        "eventType": "windowSize",
        "width": w,
        "height": h
      });
}

function post(data, callback){
    if (typeof(data)!=='object'){
        console.error("Data not recognized:",data);
    }
    data.websiteUrl = $(location).attr('href');
    data.sessionId = sessionId;
    $.ajax({
        url:"/",
        type:"POST",
        data:JSON.stringify(data),
        contentType:"application/json; charset=utf-8",
    }).done(function(data) {
        console.log(data);
    }).fail(function(xhr, status, error){
        let errorMessage = xhr.status + ': ' + xhr.statusText
        console.error(errorMessage);
    });
}
