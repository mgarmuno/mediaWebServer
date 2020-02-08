function searchMovieByTitle(e) {
    e.preventDefault();

    var xhttp = new XMLHttpRequest();

    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var movies = getMovies(this.responseText);
            showMoviesPopup(movies);
        }
    };

    var title = document.forms['movieSearchForm']['title'].value
    var page = 100;
    var year = document.forms['movieSearchForm']['year'].value;

    var data = {title: title};
    if (year) {
        data.year = year;
    }

    xhttp.open("GET", '/api/omdb/', true);
    xhttp.send(JSON.stringify(data));
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
