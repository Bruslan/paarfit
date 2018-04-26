var idleTime = 0;
var xmlhttp;
$(document).ready(function () {
    //Increment the idle time counter every minute.
    var idleInterval = setInterval(timerIncrement, 60000); // 1 minute

    //Zero the idle timer on mouse movement.
    $(this).mousemove(function (e) {
        idleTime = 0;
    });
    $(this).keypress(function (e) {
        idleTime = 0;
    });

    // Get the button, and when the user clicks on it, execute myFunction
    document.getElementById("button1").onclick = function() {myFunction()};

    /* myFunction toggles between adding and removing the show class, which is used to hide and show the dropdown content */
    function myFunction() {
        document.getElementById("button1").style.color = "red";


        xmlhttp = new XMLHttpRequest();
        xmlhttp.open("GET","/bruse", true);
        xmlhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            console.log(this.responseText)
            // var myArr = JSON.parse(this.responseText);
            // console.log(myArr)
        }
        };
        xmlhttp.send();

    }
});

function timerIncrement() {
    idleTime = idleTime + 1;
    if (idleTime > 10) { // 20 minutes
        xmlhttp = new XMLHttpRequest();
        xmlhttp.open("GET","/logout", true);
        xmlhttp.onreadystatechange = function() {
          if (this.readyState == 4 && this.status == 200) {
            window.location.reload();
          }
        };
        xmlhttp.send();
    }
}

