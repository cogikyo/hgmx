import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector("circle");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const mouth = el.querySelector('path[d^="M16.472"]');
  animate(mouth, { scaleY: [-1, 1.15, 1] }, { duration: τmid, delay: τmin });

  const eyes = el.querySelectorAll('path[d$="V9"]');
  animate(eyes, { translateY: [1, -1, 0] }, { duration: τmin, delay: τmid });
}
