function createSideNav(logoText,logoIconName,sidebarLinks){
    const sidebar = document.createElement('ul')
    sidebar.classList.add('sidenav')
    sidebar.id ="sidenav-left"
    sidebar.style ="transform: translateX(0px);"

    const logoContainer = document.createElement('li')
    logoContainer.classList.add('logo-container')
    logoContainer.classList.add('valign-wrapper')
    logoContainer.style.display = 'flex'
    logoContainer.style.alignItems = 'center'
    logoContainer.style.gap = '10px'

    const logoIcon = document.createElement('i')
    logoIcon.classList.add('material-icons','left')
    logoIcon.textContent = logoIconName

    logoContainer.appendChild(logoIcon)

    const logoTextNode = document.createTextNode(logoText)
    logoContainer.appendChild(logoTextNode)

    sidebar.appendChild(logoContainer)
    document.querySelector('header').appendChild(sidebar)

    sidebarLinks.forEach(link=>{
        const navItem = document.createElement('li')
        navItem.classList.add('bold')
        if(window.location.pathname == link.href){
            navItem.classList.add('active')
        }

        const itemLink = document.createElement('a')
        itemLink.classList.add('waves-effect')
        itemLink.href = link.href
        itemLink.textContent = link.label

        navItem.appendChild(itemLink)
        sidebar.appendChild(navItem)
    })
    // Инициализируем Materialize после создания меню
    setTimeout(() => {
        // Проверяем, что Materialize загрузился
        if (typeof M !== 'undefined' && M.Sidenav) {
            // Инициализируем sidenav
            const sidenavElems = document.querySelectorAll('.sidenav');
            M.Sidenav.init(sidenavElems, {});
            
            // Принудительно применяем стили к valign-wrapper
            const valignElems = document.querySelectorAll('.valign-wrapper');
            valignElems.forEach(el => {
                el.style.display = 'flex';
                el.style.alignItems = 'center';
                el.style.gap = '10px';
            });
        } else {
            // Если Materialize не загрузился, используем чистый CSS
            const valignElems = document.querySelectorAll('.valign-wrapper');
            valignElems.forEach(el => {
                el.style.cssText = `
                    display: flex !important;
                    align-items: center !important;
                    gap: 10px !important;
                `;
            });
        }
    }, 100);
}