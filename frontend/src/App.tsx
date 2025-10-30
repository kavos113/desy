import './App.css';
import Title from './components/Title';
import Content from './components/Content';

const App = () => {
  return (
    <main className="container">
      <div className="main">
        <Title />
        <Content />
      </div>
    </main>
  );
};

export default App;
