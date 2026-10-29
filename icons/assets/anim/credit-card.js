import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g");
  animate(card, { translateX: [0, 4, -2, 0] }, { duration: τmax });

  const chip = el.querySelector('path[d="M6 14h2"]');
  animate(chip, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
