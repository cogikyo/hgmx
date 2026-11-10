import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const bolt = el.querySelector('path[d^="m12 10"]');
  animate(bolt, { translateY: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
