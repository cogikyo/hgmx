import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M10 15"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const shield = el.querySelector('path[d^="M22 17.5"]');
  animate(shield, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
