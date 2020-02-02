function searchMovieByTitle(e) {
    e.preventDefault();

    var xhttp = new XMLHttpRequest();

    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var movies = getMovies(this.responseText);
            presentMoviesPopup(movies);
        }
    }

    var title = document.forms['movieSearchForm']['title'].value
    var year = document.forms['movieSearchForm']['year'].value;

    var data = {title: title};
    if (year) {
        data.year = year;
    }

    xhttp.open("POST", '/api/omdb/', true);
    xhttp.send(JSON.stringify(data));
}

function getMovies(json) {
    if (typeof json === 'string') {
        var movies = JSON.parse(json);
    }
    return movies;
}



// function searchMovieByTitle() {
//     var urlData = getSearchMovieData();
//     var omdbURL = '/omdbService/';
//     $.ajax({
//         url: omdbURL,
//         // context: document.body,
//         data: urlData,
//         type: "GET",
//         timeout: 4000,
//         dataType: 'json',
//         success: function(data) {
//             console.log(data);
//         },
//         error: function(error) {
//             console.log(error);
//         },
//         complete: function(complete) {
//             console.log(complete);
//         }
//     });
// }

// function getSearchMovieData() {
//     var apiKey = 'daee70b3';
//     var title = document.forms['movieSearchForm']['title'].value;
//     var year = document.forms['movieSearchForm']['year'].value;
//     var data = {
//         apikey: apiKey,
//          s: title
//     };
//     if (year) {
//         data.y = year;
//     }
//     return data;
// }

// function ajaxSuccess(data) {
//     var weba = data;
// }

// function ajaxError(error) {
//     console.log('Error in ajax call' + error.message);
// }