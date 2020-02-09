function searchMovieByTitle(e) {
    e.preventDefault();
    var title = document.forms['movieSearchForm']['title'].value
    $.ajax({
        url: 'http://localhost:8080/api/omdb/',
        data: {title: title},
        typeData: 'application/json',
        success: showMoviesPopup
    })
}

function getMovies(json) {
    if (typeof json === 'string') {
        var movies = JSON.parse(json);
    }
    return movies;
}

function uploadMovies(e, movies) {
    e.preventDefault();

    var xhttp = new XMLHttpRequest();

    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            console.log(this.responseText)
        }
    };

    var formData = new FormData();
    formData.append('file', movies[0]);

    xhttp.open("POST", '/api/file/', true);
    xhttp.send(formData);
}
