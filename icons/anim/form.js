import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const fields = el.querySelectorAll("rect");
  animate(fields, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const labels = [el.querySelector('path[d="M4 2h10"]'), el.querySelector('path[d="M4 14h6"]')];
  animate(labels, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
