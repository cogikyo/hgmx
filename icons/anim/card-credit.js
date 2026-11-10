import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g");
  animate(card, { translateX: [0, 4, -2, 0] }, { duration: τmax });

  const stripe = el.querySelector('path[d="M2 10h20"]');
  animate(stripe, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
