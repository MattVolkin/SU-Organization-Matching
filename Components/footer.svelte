<div class="footer">

  <div class="capstone-info">
    <p>This tool was made as Computer Science Capstone at Southwestern University<br>By: Aidan Balakrishnan 26’, Tanner Klein 26’, Ben McKallip 26’ and Matthew Volkin 26’</p>
    <p></p>
  </div>
</div>

<style>
  .footer {
    position: relative;
    width: 100%;
    box-sizing: border-box;
    background-color: #1f2937;
    color: #fff;
    padding: 1rem;
    box-shadow: 0 -2px 6px rgba(0, 0, 0, 0.2);
    margin-top: auto;
    text-align: center;
  }

  .capstone-info {
    overflow-wrap: anywhere;
  }

  .capstone-info p {
    margin: 0.5rem 0;
    font-size: 0.875rem;
    line-height: 1.5;
  }

  /* Desktop styles */
  @media (min-width: 768px) {
    .footer {
      padding: 0.75rem 1rem;
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 2rem;
      text-align: left;
    }

    .capstone-info {
      flex: 1 1 320px;
      text-align: center;
    }

    .capstone-info p {
      font-size: 1rem;
    }
  }
</style>