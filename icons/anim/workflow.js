import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const link = el.querySelector("path");
  animate(link, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const target = el.querySelector('rect[x="13"]');
  animate(target, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
