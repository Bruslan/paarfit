'use strict';

//////////////////////////////////////////////////
/*Sign up functionality check*/
var verifyCallback = function(response) {
  alert(response);
};

var onloadCallback = function() {
  grecaptcha.render('captcha_element', {
    'sitekey' : 'your_site_key'
  });
};

// verify sign up:
$('#sign-up').on('submit', function( event ) {

  // prevent reload
  event.preventDefault();

  // clears appended msg
  clearMsg();

  // data parsing:
  var params = {};
  params['first_name'] = document.getElementsByName('first_name')[0].value;
  params['last_name'] = document.getElementsByName('last_name')[0].value;
  params['email'] = document.getElementsByName('email2')[0].value;
  params['company'] = document.getElementsByName('company')[0].value;
  params['country'] = document.getElementsByName('country')[0].value;  
  params['passw1'] = document.getElementsByName('passw1')[0].value;
  params['passw2'] = document.getElementsByName('passw2')[0].value;
  params['age_month'] = document.getElementsByName('age_month')[0].value;
  params['age_day'] = document.getElementsByName('age_day')[0].value;
  params['age_year'] = document.getElementsByName('age_year')[0].value;

  // check pw:
  if (params['passw1'].length < 8 || params['passw1'].length > 20) {
    addMsg("Your password must be 8-20 characters long.");
    return;
  } 
  if (!validPW(params['passw1'])) {
    addMsg("Your password must contain letters, capital letters, and numbers, and must not contain spaces, special characters, or emoji.");
    return;
  }
  if (params['passw1'] != params['passw2']) {
    addMsg("Your passwords must match.");
    return;
  }
  // verify age:
  var user_date = new Date(parseInt(params['age_year']), parseInt(params['age_month']), parseInt(params['age_day']));
  var ageDif = new Date(Date.now() - user_date.getTime());
  if (Math.abs(ageDif.getUTCFullYear() - 1970) < 16) {
    addMsg("You must be at least 16 years old.");
    return;
  }

  // insert parameters into database:
  var r = new XMLHttpRequest();
  r.open("POST", "/signup_account", true);
  /*r.responseType = 'text';*/
  r.setRequestHeader("Content-Type", "application/json");
  r.onreadystatechange = function() {
      if (r.readyState === 4 && r.status === 200) {
        var json_resp = JSON.parse(r.responseText);
        if (json_resp["success"] == "true") {
          clearMsg();
          addMsg(json_resp["msg"] + " You can now login.");
          addMsg("Page reload in 5 seconds.")
          setTimeout(function(){
            location.reload();
          }, 5000);
        } else {
          clearMsg();
          addMsg(json_resp["msg"])
        }
        return false;
      }
  };
  r.send(JSON.stringify(params));
});

function validPW(str_in) {
  if (/[0-9]/.test(str_in) && /[a-z]/.test(str_in) && /[A-Z]/.test(str_in) && /^[a-zA-Z0-9]+$/.test(str_in)) {
    return true
  } else {
    return false
  }
}

// function to append div text to sign up:
function addMsg(str_in) {
  var div = document.createElement('div');
  div.className = 'container text-center';
  div.innerHTML = str_in;
  document.getElementById('responsemsg').appendChild(div);
}

// clears all appended info in lowest div:
function clearMsg() {
  var div = document.getElementById("responsemsg");
  while (div.firstChild) {
      div.removeChild(div.firstChild);
  }
}


