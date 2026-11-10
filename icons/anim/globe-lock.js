import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const globe = el.querySelectorAll('path[d^="M15.686"], path[d="M2 12h8.5"]');
  animate(globe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const shackle = el.querySelector('path[d^="M20 6V4"]');
  animate(shackle, { translateY: [0, 1.5, 0] }, { duration: τmid, delay: τmin });
}
