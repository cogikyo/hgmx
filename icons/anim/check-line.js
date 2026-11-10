import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const arms = [el.querySelector('path[d^="M20 4"]'), el.querySelector('path[d^="M9 15"]')];
  animate(arms, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(τmin) });
  const tick = el.querySelector("g");
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
  const line = el.querySelector('path[d^="M21 19"]');
  animate(line, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}
