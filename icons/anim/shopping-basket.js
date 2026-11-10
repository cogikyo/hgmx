import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const basket = el.querySelectorAll('path[d^="M2 11"], path[d^="m3.5 11"], path[d^="M4.5 15.5"], path[d^="m9 11"], path[d^="m15 11"]');
  animate(basket, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const handles = el.querySelectorAll('path[d^="m5 11"], path[d^="m19 11"]');
  animate(handles, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
