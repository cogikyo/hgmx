import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stand = el.querySelector('line[x1="12"]');
  animate(stand, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const base = el.querySelector('line[x1="8"]');
  animate(base, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
