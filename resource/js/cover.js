//<script>


function myFunction() {
    var video = document.getElementById("myVideo");
    var btn = document.getElementById("myBtn");
    if (video.paused) {
        video.play();
        btn.innerHTML = "Pause";
    } else {
        video.pause();
        btn.innerHTML = "Play";
    }
}
//</script>

//<script>
// Get the modal


// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    var modal = document.getElementById('id01');
    if (event.target == modal) {
        modal.style.display = "none";
    }

    var modal = document.getElementById('id02');
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

// </script>