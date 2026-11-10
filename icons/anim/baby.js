import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector('path[d^="M19.38"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mouth = el.querySelector('path[d^="M10 16"]');
  animate(mouth, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
