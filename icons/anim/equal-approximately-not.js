import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = [el.querySelector('path[d^="M5 9"]'), el.querySelector('path[d^="M5 15"]')];
  animate(waves, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12) });

  const slash = el.querySelector("line");
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
