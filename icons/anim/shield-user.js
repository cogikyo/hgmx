import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shield = el.querySelector('path[d^="M20 13"]');
  animate(shield, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const body = el.querySelector('path[d^="M6.376"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
