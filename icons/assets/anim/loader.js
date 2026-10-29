import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const wheel = el.querySelector("g");
  animate(wheel, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
  const spokes = el.querySelectorAll("path");
  animate(spokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
