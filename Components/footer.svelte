<div class="footer">
  <div class="content-links">
    <p>Links</p>
    <ul>
      <li><a href="/">Home</a></li>
      <li><a href="/about">About this Project</a></li>
      <li><a href="/howto">How To Use this tool</a></li>
    </ul>
  </div>

  <div class="capstone-info">
    <p>This tool was made as Computer Science Capstone at Southwestern University</p>
    <p>By: Aidan Balakrishnan 26’, Tanner Klein 26’, Ben McKallip 26’ and Matthew Volkin 26’</p>
  </div>
</div>

<style>
  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    background-color: #1f2937;
    color: #fff;
    padding: 0.75rem 1rem;
    box-shadow: 0 -2px 6px rgba(0, 0, 0, 0.2);
    z-index: 1000;

    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 2rem;
  }

  .content-links {
    flex: 1;
    text-align: left;
  }

  .capstone-info {
    flex: 1;
    text-align: right;
  }

  .content-links ul {
    margin: 0;
    padding-left: 1rem;
  }

  .content-links a {
    color: #fff;
    text-decoration: none;
  }

  .content-links a:hover {
    color: #93c5fd;
  }
</style>