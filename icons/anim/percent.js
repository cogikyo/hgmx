import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const slash = el.querySelector("line");
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const dots = [el.querySelector('circle[cx="6.5"]'), el.querySelector('circle[cx="17.5"]')];
  animate(dots, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
