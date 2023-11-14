function setIcon(isDark) {
    const themeIcon = document.getElementById('themeIcon');
    themeIcon.textContent = isDark ? '☾' : '☼';
}


function toggleTheme() {
    const htmlEl = document.documentElement;
    const currentTheme = htmlEl.getAttribute('data-bs-theme');

    if (currentTheme === 'dark') {
        htmlEl.setAttribute('data-bs-theme', 'light');
        localStorage.setItem('theme', 'light');
        setIcon(false);
    } else {
        htmlEl.setAttribute('data-bs-theme', 'dark');
        localStorage.setItem('theme', 'dark');
        setIcon(true);
    }
}

function loadTheme() {
    const storedTheme = localStorage.getItem('theme');
    const darkModeSwitch = document.getElementById('darkModeSwitch');

    if (storedTheme) {
        document.documentElement.setAttribute('data-bs-theme', storedTheme);
        darkModeSwitch.checked = storedTheme === 'dark';
        setIcon(storedTheme === 'dark');
    }  else {
        setIcon(false);
    }
}

document.getElementById('darkModeSwitch').addEventListener('click', toggleTheme);

// Call loadTheme when the page loads
loadTheme();
