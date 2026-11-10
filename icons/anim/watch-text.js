import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const lines = [el.querySelector('path[d="M10 10h2"]'), el.querySelector('path[d="M14 14h-4"]')];
  animate(lines, { opacity: [0, 1], x: [-2, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
