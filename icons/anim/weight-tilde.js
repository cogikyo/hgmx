import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const weight = el.querySelector("g");
  animate(weight, { translateY: [0, -3, 0] }, { duration: τmid });
  animate(weight, { scaleY: [1, 0.9, 1] }, { duration: τmin, delay: τmid });

  const tilde = el.querySelector('path[d^="M7.999"]');
  animate(tilde, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
