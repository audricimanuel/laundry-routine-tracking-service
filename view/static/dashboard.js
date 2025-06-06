document.addEventListener("DOMContentLoaded", () => {
  const buttons = document.querySelectorAll("button.view-details");

  buttons.forEach(button => {
    button.addEventListener("click", () => {
      const card = button.closest("div.bg-white");
      const title = card.querySelector("h3").textContent.trim();
      console.log("Clicked View Details for:", title);

      document.getElementById('viewDetailsModal').classList.remove('hidden');
      document.getElementById('viewDetailsModal').classList.add('flex');
    });
  });

  const closeViewDetailsModalButton = document.getElementById("closeViewDetailsModal")

  closeViewDetailsModalButton.addEventListener("click", () => {
    document.getElementById('viewDetailsModal').classList.remove('flex');
    document.getElementById('viewDetailsModal').classList.add('hidden');
  })
});