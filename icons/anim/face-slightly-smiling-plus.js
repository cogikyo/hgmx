import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector('path[d^="M13.267"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mouth = el.querySelector('path[d^="M16.472"]');
  animate(mouth, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });

  const eyes = el.querySelectorAll('path[d$="V9"]');
  animate(eyes, { translateY: [1, -1, 0] }, { duration: τmin, delay: τmid });

  const plus = el.querySelectorAll('path[d="M16 5h6"], path[d="M19 2v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
