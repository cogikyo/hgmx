import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bars = [el.querySelector('line[y1="9"]'), el.querySelector('line[y1="15"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: stagger(0.12) });

  const slash = el.querySelector('line[y1="5"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
