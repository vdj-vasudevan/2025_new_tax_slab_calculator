document.getElementById("add-bucket").addEventListener("click", () => {
  const div = document.createElement("div");
  div.className = "flex space-x-2";
  div.innerHTML = `
    <input type="number" step="any" placeholder="From" class="w-1/3 p-2 border rounded" />
    <input type="number" step="any" placeholder="To" class="w-1/3 p-2 border rounded" />
    <input type="number" step="any" placeholder="%" class="w-1/3 p-2 border rounded" />
  `;
  document.getElementById("custom-buckets").appendChild(div);
});

document.getElementById("custom-tax-form").addEventListener("submit", async function (e) {
  e.preventDefault();

  const annualPay = parseFloat(this.annual_pay.value);
  const bucketDivs = document.querySelectorAll("#custom-buckets > div");
  const customBuckets = [];

  bucketDivs.forEach(div => {
    const inputs = div.querySelectorAll("input");
    customBuckets.push({
      from: parseFloat(inputs[0].value),
      to: parseFloat(inputs[1].value),
      percent: parseFloat(inputs[2].value)
    });
  });

  try {
    const res = await fetch("/calculate-custom-tax", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ annual_pay: annualPay, buckets: customBuckets }),
    });

    if (!res.ok) throw new Error("Failed to get response");

    const data = await res.json();
    document.getElementById("custom-tax").textContent = data.tax.toFixed(2);
    document.getElementById("custom-inhand").textContent = data.inhand.toFixed(2);
    document.getElementById("result").classList.remove("hidden");
  } catch (error) {
    alert("Custom tax calculation failed: " + error.message);
  }
});
