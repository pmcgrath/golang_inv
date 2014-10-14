package main

const rootHtmlTemplate = `
<html>
  <head>
    <title>Contacts</title>
    <script src="/assets/js/main.js"></script>    
  </head>
  <body>
    <h1>Contacts</h1>
    <div id="welcomeSection">
      Welcome <span id="welcomeUserName"></span>
      <button id="logOut">Log out</button>
    </div>
    <div id="logInSection">
      <form name="logInForm">
        <input type="text" id="userName" placeholder="username" autofocus>
        <input type="password" id="password" placeholder="password">
        <button id="logIn">Log in</button>
      </form>
      <br/>
      <span id="logInMessage"></span>
    </div>
    <div id="contactListSection">
      <h2>Contacts</h2>
      <div id="contactList"></div>
      <button id="newContact">New contact</button>
    </div>
    <!-- See http://www.html5rocks.com/en/tutorials/webcomponents/template/ -->
    <template id="contactListItemTemplate">
      <div class="contactListItem">
        <label class="contactListItemFirstName"></label> <!-- Needed to use closing tags, again chrome had a problem i do not understand at this time -->
        <label class="contactListItemLastName"></label>
        <button class="contactListItemEdit">e</button>
        <button class="contactListItemDeletion">x</button>
      </div>
    </template>
    <dialog id="contactEditor">
      <form name="contactEditorForm">
        <input type="hidden" id="contactId"/>
        <label classs="contactiEditorLabel">First name:</label><input type="text" id="contactFirstName"/><br/>
        <label classs="contactiEditorLabel">Last name:</label><input type="text" id="contactLastName"/><br/>
        <label classs="contactiEditorLabel">Emails:</label><input type="text" id="contactEmails"/><br/>
        <label classs="contactiEditorLabel">Phones:</label><input type="text" id="contactPhones"/><br/>
        <label classs="contactiEditorLabel">Twitter:</label><input type="text" id="contactTwitter" placeholder="@twitterhandle" pattern="^@?(\w){1,15}$"/><br/>
        <label classs="contactiEditorLabel">Notes:</label><input type="textarea" id="contactNotes"/><br/>
      </form>
      <!-- Need to keep buttons outside form as seems to mess up in chrome - no idea why at this time -->
      <button id="saveContact">Save</button>
      <button id="cancelEdit">Cancel</button>
    </dialog>
  </body>
</html>`
