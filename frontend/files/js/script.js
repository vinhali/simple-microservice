const body = document.querySelector(".body");
const boxContainer = document.querySelector(".box-contaner");
const boxHdr = document.querySelector(".box--hdr");
const boxInner = document.querySelector(".box--inner");

// USER INFORMATION === start
const formauth = document.querySelector(".auth");
const passwordInputInfo = document.querySelector(".passwordInput--info");
const passwordInputHdr = document.querySelector(".passwordInput");
const userInputPassword = document.querySelector(".password--input");
const userSubmit = document.querySelector(".password-Submit");
const transferData = document.querySelector(".transferData");
const data = document.querySelector(".data");
const userAvatar = document.querySelector(".user--avatar");
let userLetter = "V";

const showTranferData = function () {
  const userLetter = "V";

  if (userLetter === "" || userLetter === null) {
    userInputPassword.classList.toggle("shake-horizontal");
  } else {
    clearAllCookies();
    const userAvatar = userLetter[0].toUpperCase();
    const displayHeader = `<div class="password--info-active"><h1 class="userName--h1"><span class="userName--emoji">👋</span><span class="userName--span"> Hello, </span>Vinhali.</h1><span class="user--avatar">${userAvatar}</span></div>`;

    const dataH2 = `<div class="data-h2"><h2>Completed</span></h2></div>`;
    transferData.insertAdjacentHTML("afterbegin", dataH2);

    document.getElementById("userInfo").innerHTML = " ";
    passwordInputInfo.insertAdjacentHTML("afterbegin", displayHeader);
    boxHdr.classList.toggle("box--hdr--passActive");
    passwordInputHdr.classList.toggle("password--active");
    transferData.classList.toggle("transferData--Active");
    passwordInfoActive.classList.toggle("password--info-active");
    document.getElementById("user").innerHTML = " ";
    acctUser.classList.add("userCurrent");
    acctUser.insertAdjacentHTML("afterbegin", displayScroll);
  }
};

const returnHome = function () {
  boxHdr.classList.toggle("box--hdr--passActive");

  const isUserActive = boxHdr.classList.contains("box--hdr--passActive");

  if (isUserActive) {
    boxHdr.style.backgroundImage = 'url(https://github.com/vinhali/simple-microservice/blob/main/files/menu.jpg?raw=true&resize=800x600)';
  } else {
    boxHdr.style.backgroundImage = 'url(https://github.com/vinhali/simple-microservice/blob/main/files/menu.jpg?raw=true&resize=800x600)';
  }

  transferData.classList.remove("transferData--Active");
  passwordInputHdr.classList.remove("password--active");
  document.getElementById("userInfo").innerHTML = "";
};


const returnToMenu = document.querySelector('p[type="submit"]');

returnToMenu.addEventListener('click', () => {
  returnHome();
});


const performTransaction = async () => {
  const passwordInput = document.getElementById('password');
  const password = passwordInput.value;

  if (password === '1234') {
    try {
      const headers = new Headers();

      if (document.cookie.indexOf('authpass=ok') === -1) {
        document.cookie = 'authpass=passed';

        const response = await fetch('/auth', {
          method: 'POST',
          headers: headers,
        });

        if (!response.ok) {
          throw new Error(`Authentication failed! Status: ${response.status}`);
        }
        showTranferData();
      }

    } catch (error) {
      console.error(error);
    }
  } else {
    document.cookie = 'authpass=rejected';

    window.location.reload();
  }
};

const clearAllCookies = () => {
  const cookies = document.cookie.split(";");

  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i];
    const eqPos = cookie.indexOf("=");
    const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
    document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
  }
};

window.onload = clearAllCookies;

userSubmit.addEventListener('click', () => {
  performTransaction();
});
