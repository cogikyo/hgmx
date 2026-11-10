import { animate, τmid } from "../motion.js";

export default function (el) {
  const shaft = el.querySelector('path[d^="m13.41"]');
  animate(shaft, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
}
