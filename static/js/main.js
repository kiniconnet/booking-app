  // Toggle Navbar on Small Screens
  const navbarToggler = document.getElementById('navbarToggler');
  const navbarNav = document.getElementById('navbarNav');
  navbarToggler.addEventListener('click', () => {
      navbarNav.classList.toggle('active');
  });

  // Dropdown Menu
  const roomsDropdown = document.getElementById('roomsDropdown');
  const roomsMenu = document.getElementById('roomsMenu');
  roomsDropdown.addEventListener('click', (e) => {
      e.preventDefault();
      roomsMenu.style.display = roomsMenu.style.display === 'block' ? 'none' : 'block';
  });

  // Close Dropdown When Clicking Outside
  document.addEventListener('click', (e) => {
      if (!roomsDropdown.contains(e.target)) {
          roomsMenu.style.display = 'none';
      }
  });

 