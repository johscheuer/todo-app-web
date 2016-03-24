$(document).ready(function() {
  var entryContentElement = $("#todo-input");

  var appendTodoList = function(data) {
    if (data == null) {
      return
    }
    $("#Todos > tbody").empty();
    $.each(data, function(key, val) {
      var myRow = '<tr><td class="col-md-10">'+val+'</td><td class="vcenter col-md-2"><input type="checkbox" name="deleteCheck" value="1"/></td></tr>';
      $("#Todos > tbody").append(myRow);
    });
  }

  var handleSubmission = function(e) {
    e.preventDefault();
    var entryValue = entryContentElement.val()
    if (!entryValue || entryValue.length <= 0){
        entryContentElement.parent().addClass("has-error").removeClass("has-success");
        return false;
    }


    entryContentElement.val("")
    entryContentElement.parent().removeClass("has-error").addClass("has-success");
    $.getJSON("insert/todo/" + entryValue, appendTodoList);
  }

  var handleDeletion = function(e){
    e.preventDefault();

    var checkboxes = document.getElementsByName("deleteCheck");
    for (var i=0; i < checkboxes.length; i++) {
     if (!checkboxes[i].checked){
       continue
     }
     var checkbox = checkboxes[i];
     $.getJSON("delete/todo/" + $(checkbox).closest('tr').text(), appendTodoList);
    }
  }

  $("#todo-submit").click(handleSubmission);
  $("#todo-delete").click(handleDeletion);

  // Poll every second.
  (function fetchTodos() {
    $.getJSON("read/todo").done(appendTodoList).always(
      function() {
        setTimeout(fetchTodos, 10000);
      });
  })();
});
