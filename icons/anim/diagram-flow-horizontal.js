import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const links = el.querySelectorAll("path");
  animate(links, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const targets = [el.querySelector('rect[y="2"]'), el.querySelector('rect[y="16"]')];
  animate(targets, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
