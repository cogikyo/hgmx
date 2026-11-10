import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelectorAll('path[d^="M10.1"], path[d^="M17.1"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pacifier = el.querySelector("g");
  animate(pacifier, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
