import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const arc = el.querySelector('path[d^="M14 12"]');
  animate(arc, { rotate: [0, 720] }, { duration: τmax });
  animate(arc, { pathLength: [0.25, 1] }, { duration: τmax });
}
