import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const ribbon = el.querySelector("g");
  animate(ribbon, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmax });

  const outline = el.querySelector("path");
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
