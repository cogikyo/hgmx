import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelectorAll('path[d^="M19 11"], path[d^="M12 18"], path[d^="M9 22"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const egg = el.querySelector('path[d^="M18 11"]');
  animate(egg, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
