import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const leads = el.querySelectorAll("path");
  animate(leads, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const body = el.querySelector("rect");
  animate(body, { scale: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
