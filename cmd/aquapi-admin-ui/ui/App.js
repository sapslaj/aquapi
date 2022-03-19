import Form from './pages/Form.js';
import Random from './pages/Random.js';

const { createElement: e, memo, useState, Fragment } = React;

export default function App() {
  const [tab, setTab] = useState('Random');

  const NavItem = ({ name }) => {
    const className = tab === name ? 'nav-link active' : 'nav-link';
    const onClick = () => setTab(name);
    return e('a', { href: '#', className, onClick}, name);
  }

  const navItem = (name) => (
    e(NavItem, {key: name, name})
  )

  const PageWrapper = ({ name, component }) => {
    const display = tab === name ? 'block': 'none';
    const c = e(memo(component))
    return e('div', {style: {display}}, c);
  }

  const pageWrapper = (component) => (
    e(PageWrapper, {key: component.name, name: component.name, component: component})
  )

  return e(Fragment, null, [
    e('nav', {className: 'navbar navbar-expand navbar-dark bg-dark'},
      e('div', {className: 'container-fluid justify-content-start'}, [
        e('a', {key: 'brand', className: 'navbar-brand', href: '#'}, 'AquaPI Admin'),
        e('div', {className: 'navbar-nav'}, [
          navItem('Random'),
          navItem('Form'),
        ])
      ])
    ),
    pageWrapper(Random),
    pageWrapper(Form),
  ])
}
