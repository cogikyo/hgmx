import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const grid = el.querySelectorAll("path");
  animate(grid, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06) });
}
