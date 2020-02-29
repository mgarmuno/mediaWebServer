function onDocumentReady() {
    $('#addMovieButton').on('click', triggerOpenDialog)
}

function openPopup(popup) {
    closeAllPopups();
    $(popup).show();
}

function closeAllPopups() {
    $("#movieSearchDiv").hide();
    $('#movieSearchResultPopup').hide();
    $('#movieOptionsPopup').hide();
}

function openMovieSearchForm() {
    $("#movieSearchDiv").first().show();
}

function closeMovieSearchForm() {
}

function openMovieCardsDiv() {
    $("#moviesSearchResultPopup").first().show();
}

function triggerOpenDialog() {
    $('#fileInput').trigger('click');
}
