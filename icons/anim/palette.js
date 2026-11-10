import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const path = el.querySelector("path");
  animate(path, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const circles = el.querySelectorAll("circle");
  animate(circles, { scale: [0, 1.75, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const group = el.querySelector("g");
  animate(group, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
