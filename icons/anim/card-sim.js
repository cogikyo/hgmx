import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector('path[d^="M14.172"]');
  animate(card, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const chip = el.querySelector("g.chip");
  animate(chip, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
