import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bot = ["M13.67", "M8 8H6", "M9.67"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(bot, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="M22 22"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
