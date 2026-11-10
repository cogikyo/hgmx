import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelectorAll('path[d^="m18.2"], path[d^="m7.2"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const sleeve = el.querySelector('path[d^="M18 18"]');
  animate(sleeve, { translateY: [3, -0.5, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
