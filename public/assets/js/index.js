function getEffect() {
  console.log("getting effect");
  $.get('/magical-effect', function(data) {
    data = JSON.parse(data);
    $('.magical-effect .roll').html(data['key']);
    $('.magical-effect .text').html(data['text'])
  });
}

$(function() {
  getEffect();

  $('.new-magical-effect').click(function() {
    getEffect();
  });
});

